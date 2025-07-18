const puppeteer = require('puppeteer')

let browserInstance = null

const delay = ms => new Promise(resolve => setTimeout(resolve, ms))

const getBrowser = async () => {
	if (!browserInstance) {
		browserInstance = await puppeteer.launch({
			headless: true,
			args: [
				'--no-sandbox',
				'--disable-setuid-sandbox',
				'--disable-dev-shm-usage',
				'--disable-gpu',
				'--disable-extensions',
				'--disable-plugins',
				'--disable-default-apps',
				'--no-first-run',
			],
			defaultViewport: { width: 1280, height: 720 },
		})
	}
	return browserInstance
}

const closeBrowser = async () => {
	if (browserInstance) {
		await browserInstance.close()
		browserInstance = null
	}
}

const scrapeAdvertiser = async domain => {
	const browser = await getBrowser()
	const page = await browser.newPage()

  console.time(`[⏱ ${domain}]`)

	try {
		// Отключаем загрузку изображений
		await page.setRequestInterception(true)
		page.on('request', req => {
			const type = req.resourceType()
			if (type === 'image' || type === 'media') req.abort()
			else req.continue()
		})

		await page.setUserAgent(
			'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36'
		)

		await page.goto('https://adstransparency.google.com/?region=KZ', {
			waitUntil: 'networkidle0',
			timeout: 30000,
		})

		await page.waitForSelector('input', { timeout: 10000 })
		await page.type('input', domain, { delay: 50 })
		await delay(2000)

		// навигация по подсказке
		await page.keyboard.press('Tab')
		await page.keyboard.press('Tab')
		await page.keyboard.press('Enter')
		await delay(500)
		await page.keyboard.press('Tab')
		await page.keyboard.press('Tab')
		await page.keyboard.press('Enter')

		// ждем блок с данными
		await page.waitForSelector('.metadata', { timeout: 10000 })

		const result = await page.evaluate(() => {
			const metadata = document.querySelector('.metadata')
			if (!metadata) return { legal_name: '', country: '' }

			const result = { legal_name: '', country: '' }
			const strongs = metadata.querySelectorAll('strong')

			for (const strong of strongs) {
				const label = strong.textContent?.trim()
				const value = strong.parentElement?.textContent
					?.replace(label, '')
					.trim()

				if (label.includes('Юридическое название')) result.legal_name = value
				if (label.includes('Страна регистрации')) result.country = value

				if (result.legal_name && result.country) break
			}

			return result
		})

    console.timeEnd(`[⏱ ${domain}]`)
		return result
	} catch (e) {
    console.timeEnd(`[⏱ ${domain}]`)
		console.error('Ошибка скрейпинга:', e.message)
		return {
			domain,
			error: '❌ Не удалось спарсить данные с сайта.',
		}
	} finally {
		await page.close()
	}
}

const scrapeMultipleDomains = async domains => {
	const results = []
	for (const domain of domains) {
		console.log(`🔍 Проверяем: ${domain}`)
		const res = await scrapeAdvertiser(domain)
		results.push({ domain, ...res })
		await delay(500)
	}
	return results
}

process.on('SIGINT', async () => {
	console.log('🛑 Завершение...')
	await closeBrowser()
	process.exit(0)
})

process.on('SIGTERM', async () => {
	console.log('🛑 Завершение...')
	await closeBrowser()
	process.exit(0)
})

// тест
const test = async () => {
	const domains = ['trust-artamonov.com', 'google.com']
	const results = await scrapeMultipleDomains(domains)
	console.log('✅ Результаты:', results)
	await closeBrowser()
}

test()
