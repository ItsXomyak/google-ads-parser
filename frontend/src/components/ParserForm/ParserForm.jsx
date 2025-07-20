import React, { useState } from 'react'
import './ParserForm.css'

const ParserTester = () => {
	const [domains, setDomains] = useState(
		'trust-artamonov.com\ngoogle.com\nfacebook.com'
	)
	const [results, setResults] = useState([])
	const [loading, setLoading] = useState(false)
	const [error, setError] = useState('')
	const [copiedIndex, setCopiedIndex] = useState(null)

	const handleParse = async () => {
		setLoading(true)
		setError('')
		setResults([])

		try {
			const domainList = domains.split('\n').filter(d => d.trim())

			const response = await fetch('http://localhost:8080/parse/batch', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify({ domains: domainList }),
			})

			if (!response.ok) {
				throw new Error(`HTTP ${response.status}: ${response.statusText}`)
			}

			const data = await response.json()
			setResults(data)
		} catch (err) {
			setError(err.message)
		} finally {
			setLoading(false)
		}
	}

	const clearResults = () => {
		setResults([])
		setError('')
	}

	const clearAll = () => {
		setResults([])
		setError('')
		setDomains('')
	}

	const copyResult = async (result, index) => {
		let textToCopy = `Домен: ${result.domain}\n`

		if (result.error) {
			textToCopy += `Ошибка: ${result.error}`
		} else {
			textToCopy += `Юридическое название: ${
				result.legal_name || 'Не найдено'
			}\n`
			textToCopy += `Страна: ${result.country || 'Не найдено'}`
		}

		try {
			await navigator.clipboard.writeText(textToCopy)
			setCopiedIndex(index)
			setTimeout(() => setCopiedIndex(null), 2000)
		} catch (err) {
			console.error('Ошибка копирования:', err)
		}
	}

	const copyAllResults = async () => {
		let allText = 'Результаты парсинга:\n\n'

		results.forEach((result, index) => {
			allText += `${index + 1}. Домен: ${result.domain}\n`
			if (result.error) {
				allText += `   Ошибка: ${result.error}\n\n`
			} else {
				allText += `   Юридическое название: ${
					result.legal_name || 'Не найдено'
				}\n`
				allText += `   Страна: ${result.country || 'Не найдено'}\n\n`
			}
		})

		try {
			await navigator.clipboard.writeText(allText)
			setCopiedIndex('all')
			setTimeout(() => setCopiedIndex(null), 2000)
		} catch (err) {
			console.error('Ошибка копирования:', err)
		}
	}

	return (
		<div className='parser-container'>
			<h2>Parser Tester</h2>

			<div className='input-section'>
				<label>Домены (каждый с новой строки):</label>
				<textarea
					value={domains}
					onChange={e => setDomains(e.target.value)}
					placeholder='example.com'
					rows={5}
					disabled={loading}
				/>

				<div className='button-group'>
					<button
						onClick={handleParse}
						disabled={loading || !domains.trim()}
						className={`parse-btn ${loading ? 'loading' : ''}`}
					>
						{loading ? 'Парсинг...' : 'Запустить парсер'}
					</button>

					<button
						onClick={clearAll}
						disabled={loading}
						className='clear-btn'
						title='Очистить всё'
					>
						Очистить всё
					</button>
				</div>
			</div>

			{error && <div className='error'>Ошибка: {error}</div>}

			{results.length > 0 && (
				<div className='results'>
					<div className='results-header'>
						<h3>Результаты:</h3>
						<div className='results-actions'>
							<button
								onClick={handleParse}
								disabled={loading || !domains.trim()}
								className='retry-btn'
								title='Повторить запрос'
							>
								🔄 Повторить
							</button>
							<button
								onClick={copyAllResults}
								className='copy-all-btn'
								title='Скопировать все результаты'
							>
								{copiedIndex === 'all' ? '✅ Скопировано' : '📋 Копировать всё'}
							</button>
							<button
								onClick={clearResults}
								className='clear-results-btn'
								title='Очистить результаты'
							>
								🗑 Очистить
							</button>
						</div>
					</div>

					{results.map((result, index) => (
						<div
							key={index}
							className={`result-item ${
								result.error ? 'error-result' : 'success-result'
							}`}
						>
							<div className='result-header'>
								<h4>{result.domain}</h4>
								<button
									onClick={() => copyResult(result, index)}
									className='copy-btn'
									title='Скопировать результат'
								>
									{copiedIndex === index ? '✅' : '📋'}
								</button>
							</div>

							{result.error ? (
								<p className='error-text'>Ошибка: {result.error}</p>
							) : (
								<div className='success-data'>
									<p>
										<strong>Юридическое название:</strong>{' '}
										{result.legal_name || 'Не найдено'}
									</p>
									<p>
										<strong>Страна:</strong> {result.country || 'Не найдено'}
									</p>
								</div>
							)}
						</div>
					))}
				</div>
			)}
		</div>
	)
}

export default ParserTester
