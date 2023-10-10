const inLangSelect = document.querySelector('#input-lang');
const outLangSelect = document.querySelector('#output-lang');
const langSwapper = document.querySelector('#swap-langs');

const inTextArea = document.querySelector('#input-text');
const outTextArea = document.querySelector('#output-text');

const translateButton = document.querySelector('#submit-translation');

langSwapper.addEventListener('click', (e) => {
    const inLang = inLangSelect.value;
    const outLang = outLangSelect.value;

    inLangSelect.value = outLang;
    outLangSelect.value = inLang;
});

translateButton.addEventListener('click', (e) => {
    console.log('translate');
});
