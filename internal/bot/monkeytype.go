package bot

import (
	"fmt"
	"time"

	"typingmonkey/internal/input"

	"github.com/tebeka/selenium"
)

type MonkeyTypeBot struct {
	driver selenium.WebDriver
	typer  *input.Typer
}

func New(driver selenium.WebDriver, typingDelay time.Duration) *MonkeyTypeBot {
	return &MonkeyTypeBot{
		driver: driver,
		typer:  input.NewTyper(typingDelay),
	}
}

func (b *MonkeyTypeBot) OpenMonkeyType() error {
	if err := b.driver.Get("https://monkeytype.com/"); err != nil {
		return fmt.Errorf("page failed to load :( oopsies!: %w", err)
	}

	time.Sleep(3 * time.Second)

	rejectBtn, err := b.driver.FindElement(selenium.ByClassName, "rejectAll")
	if err == nil {
		rejectBtn.Click()
		time.Sleep(1 * time.Second)
	}

	typingArea, err := b.driver.FindElement(selenium.ByID, "wordsWrapper")
	if err == nil {
		typingArea.Click()
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

func (b *MonkeyTypeBot) RunTest() (int, error) {
	fmt.Println("\nstarting in 3 seconds!")
	for i := 3; i > 0; i-- {
		fmt.Printf("%d...\n", i)
		time.Sleep(1 * time.Second)
	}
	fmt.Println("typing now! (▰˘◡˘▰)")

	wordsTyped := 0

	for {
		word, err := b.getActiveWord()
		if err != nil {
			break
		}

		if word == "" {
			time.Sleep(50 * time.Millisecond)
			continue
		}

		fmt.Printf("%d | %s\n", wordsTyped+1, word)

		b.typer.TypeText(word)
		b.typer.TypeSpace()

		wordsTyped++
		time.Sleep(10 * time.Millisecond)
	}

	return wordsTyped, nil
}

func (b *MonkeyTypeBot) Restart() error {
	if err := b.driver.Refresh(); err != nil {
		return err
	}

	time.Sleep(2 * time.Second)

	// Re-focus typing area
	typingArea, err := b.driver.FindElement(selenium.ByID, "wordsWrapper")
	if err == nil {
		typingArea.Click()
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

func (b *MonkeyTypeBot) getActiveWord() (string, error) {
	activeWord, err := b.driver.FindElement(selenium.ByCSSSelector, "div.word.active")
	if err != nil {
		return "", err
	}

	letters, err := activeWord.FindElements(selenium.ByTagName, "letter")
	if err != nil || len(letters) == 0 {
		return "", nil
	}

	word := ""
	for _, letter := range letters {
		text, err := letter.Text()
		if err != nil {
			continue
		}
		word += text
	}

	return word, nil
}
