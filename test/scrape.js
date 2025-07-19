const puppeeteer = require('puppeteer');

(async () => {
	const browser = await puppeeteer.launch({ headless: false }); // не показываем браузер
	const page = await browser.newPage();

  await page.goto('https://news.ycombinator.com');

  await page.waitForSelector(".titleline");

  const news = await page.evaluate(() => {
    // Извлекаем заголовки и ссылки на новости
    const items = Array.from(document.querySelectorAll('.titleline'));
    return items.map(item => { // для каждого элемента извлекаем заголовок и ссылку
      const title = item.querySelector('a')?.href || ""; // получаем ссылку на новость
      const link = item.querySelector('a')?.innerText || ""; // получаем текст заголовка
      return { title, link };
    });
  });

  console.log(news.slice(0, 5)); // выводим первые 5 новостей

	await browser.close()
})();