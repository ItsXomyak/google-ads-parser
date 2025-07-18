  const puppeteer = require('puppeteer')

  const scrapeAdvertiser = async domain => {
    const browser = await puppeteer.launch({
      headless: true,
      args: ['--no-sandbox', '--disable-setuid-sandbox'],
      defaultViewport: null,
    })

    const page = await browser.newPage()

    await page.goto('https://adstransparency.google.com/?region=KZ', {
      waitUntil: 'networkidle2',
      timeout: 0,
    })

    // Ждём поле поиска
    await page.waitForSelector('input', { timeout: 30000 })

    // Вводим домен
    await page.type('input', domain, { delay: 100 })

    // Немного ждём, чтобы подсказка успела появиться
    await new Promise(resolve => setTimeout(resolve, 2000))

    // Нажимаем TAB 2 раза и ENTER — перемещаемся на подсказку и выбираем её
    await page.keyboard.press('Tab')
    await page.keyboard.press('Tab')
    await page.keyboard.press('Enter')
    await new Promise(resolve => setTimeout(resolve, 500))
    await page.keyboard.press('Tab')
    await page.keyboard.press('Tab')
    await page.keyboard.press('Enter')

    try {
      await page.waitForSelector('.metadata', { timeout: 10000 })

      const result = await page.evaluate(() => {
        const metadata = document.querySelector('.metadata')
        if (!metadata) return { legal_name: '', country: '' }

        const strongs = metadata.querySelectorAll('strong')
        let legal_name = ''
        let country = ''

        strongs.forEach(strong => {
          const label = strong.textContent?.trim()
          const value = strong.parentElement?.textContent
            ?.replace(label, '')
            .trim()

          if (label.includes('Юридическое название')) legal_name = value
          if (label.includes('Страна регистрации')) country = value
        })

        return { legal_name, country }
      })

      return result
    } catch (e) {
      return {
        error:
          '❌ Не удалось перейти на страницу рекламодателя или спарсить данные.',
      }
    } finally {
      await browser.close()
    }
  }

  // 👇 Тест
  ;(async () => {
    const domain = 'trust-artamonov.com'
    const result = await scrapeAdvertiser(domain)
    console.log('✅ Результат:', result)
  })()
