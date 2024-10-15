package form

var FormTmpl = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Установка количества потоков для отправки данных на сервер</title>
	</head>
	<body>
			<h2>Настройка количества потоков</h2>
			<form action="/submit" method="post">
					<label for="numThreads">Количество потоков:</label>
					<input type="number" id="numThreads" name="numThreads" min="1" max="10" value="1" required>
					<br><br>
					<input type="submit" value="Установить количество потоков">
			</form>
	</body>
	</html>
`
