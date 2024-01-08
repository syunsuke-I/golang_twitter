// ツイートの内容が空もしくはスペースであれば送信ボタンを非活性化する
document.addEventListener('DOMContentLoaded', function () {
  let textarea = document.getElementById('content');
  let tweetButton = document.getElementById('submit-button');
  textarea.addEventListener('input', function () {
    tweetButton.disabled = !textarea.value.trim();
  });
});