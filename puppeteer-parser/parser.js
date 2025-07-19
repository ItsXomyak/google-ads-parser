const express = require('express')
const puppeteer = require('puppeteer')
const app = express()

app.use(express.json())

app.post('/parse-batch', async (req, res) => {
	const domains = req.body.domains
	if (!Array.isArray(domains) || domains.length === 0) {
		return res.status(400).json({ error: 'domains[] required' })
	}

	const browser = await puppeteer.launch({
		headless: true,
		args: ['--no-sandbox', '--disable-setuid-sandbox'],
	})

	try {
		const results = await Promise.all(
			domains.map(async domain => {
				const page = await browser.newPage()
				try {
					await page.goto('https://adstransparency.google.com/?region=KZ', {
						waitUntil: 'domcontentloaded',
						timeout: 10000,
					})

					await page.waitForSelector('input', { timeout: 5000 })
					await page.type('input', domain, { delay: 100 })
					await new Promise(resolve => setTimeout(resolve, 1000))
					await page.keyboard.press('Tab')
					await page.keyboard.press('Tab')
					await page.keyboard.press('Enter')
					await new Promise(resolve => setTimeout(resolve, 350))
					await page.keyboard.press('Tab')
					await page.keyboard.press('Tab')
					await page.keyboard.press('Enter')
					await page.waitForSelector('.metadata', { timeout: 5000 })

					const result = await page.evaluate(() => {
						const metadata = document.querySelector('.metadata')
						if (!metadata) return { legal_name: '', country: '' }

						const strongs = metadata.querySelectorAll('strong')
						let legal_name = '',
							country = ''
						strongs.forEach(strong => {
							const label = strong.textContent?.trim()
							const value = strong.parentElement?.textContent
								.replace(label, '')
								.trim()
							if (label.includes('Юридическое название')) legal_name = value
							if (label.includes('Страна регистрации')) country = value
						})
						return { legal_name, country }
					})

					await page.close()
					return { domain, ...result }
				} catch (e) {
					await page.close()
					return { domain, error: e.message }
				}
			})
		)

		res.json(results)
	} catch (e) {
		res.status(500).json({ error: 'batch parse failed', details: e.message })
	} finally {
		await browser.close()
	}
})

app.listen(3001, () => console.log('Parser batch server running on :3001'))
