const express = require('express');
var cookieParser = require('cookie-parser');
var bodyParser = require('body-parser');

const app = express();
const port = process.env.PORT || 8000;

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({
  extended: false
}));
app.use(cookieParser());
app.use(express.static('../src/wilson'));


app.get('/logout', (req, res) => {
  console.log('cookies', req.cookies)
  res.clearCookie('id')
  res.redirect('http://localhost:8080/');
})

app.get('/*', (req, res) => {
  res.sendStatus(400);
});

app.listen(port, () => {
  console.log(port);
});
