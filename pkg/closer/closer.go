package closer

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	// successIdx - индекс, который возвращается в случае успешного закрытия всех ресурсов
	successIdx = -1
)

// Closer обеспечивает потокобезопасное закрытие ресурсов.
type Closer struct {
	funcs         []Func
	mu            sync.Mutex
	once          sync.Once
	forcedTimeout time.Duration
}

// Func — сигнатура функции закрытия ресурса.
type Func func(ctx context.Context) error

// NewCloser создает новый экземпляр Closer.
// forcedTimeout — время, отводимое на принудительное закрытие всех ресурсов при таймауте контекста в Close.
func NewCloser(forcedTimeout time.Duration) *Closer {
	const (
		defaultForcedTimeout = 2 * time.Second
	)

	if forcedTimeout == 0 {
		forcedTimeout = defaultForcedTimeout
	}

	return &Closer{
		forcedTimeout: forcedTimeout,
	}
}

// Add добавляет функцию в список закрытия
func (c *Closer) Add(f Func) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.funcs = append(c.funcs, f)
}

// Close последовательно запускает закрытие всех зарегистрированных функций (LIFO).
// Если контекст отменяется до завершения, оставшиеся функции закрываются принудительно.
func (c *Closer) Close(ctx context.Context) error {
	var err error
	c.once.Do(func() {
		c.mu.Lock()
		funcs := c.funcs
		c.mu.Unlock()

		stopIdx, errors := c.gracefulClose(ctx, funcs)
		if stopIdx == successIdx { // Если все ресурсы закрылись успешно
			if len(errors) > 0 {
				err = fmt.Errorf("shutdown finished with error(s):\n%s", strings.Join(errors, "\n"))
			}

			return
		}

		// Если есть незакрытые ресурсы, пытаемся закрыть их принудительно
		remaining := funcs[:stopIdx+1]
		forcedErrs := c.forcedClose(remaining)
		errors = append(errors, forcedErrs...)

		err = fmt.Errorf(
			"shutdown interrupted after %d/%d funcs:\n%s",
			len(funcs)-1-stopIdx,
			len(funcs),
			strings.Join(errors, "\n"),
		)
	})

	return err
}

// gracefulClose закрывает все функции в порядке LIFO.
// Если какая-то функция возвращает ошибку, она добавляется в список ошибок.
// Если контекст будет отменен, функция вернет индекс последней успешно закрытой функции и список ошибок.
func (c *Closer) gracefulClose(ctx context.Context, funcs []Func) (int, []string) {
	var errors []string
	for i := len(funcs) - 1; i >= 0; i-- {
		var (
			f    = funcs[i]
			done = make(chan error, 1)
		)

		go func() {
			done <- f(ctx)
		}()

		select {
		case err := <-done:
			if err != nil {
				errors = append(errors, fmt.Sprintf("[!] %v", err))
			}
		case <-ctx.Done():
			return i, errors
		}
	}

	return successIdx, errors
}

// forcedClose параллельно запускает все оставшиеся функции закрытия с собственным таймаутом.
func (c *Closer) forcedClose(funcs []Func) []string {
	var (
		wg     sync.WaitGroup
		mu     sync.Mutex
		errors []string
	)

	ctx, cancel := context.WithTimeout(context.Background(), c.forcedTimeout)
	defer cancel()

	for _, f := range funcs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := f(ctx); err != nil {
				mu.Lock()
				errors = append(errors, fmt.Sprintf("[FORCED] %v", err))
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	return errors
}
