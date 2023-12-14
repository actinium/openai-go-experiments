const inLangSelect = document.querySelector('#input-lang');
const outLangSelect = document.querySelector('#output-lang');
const langSwapper = document.querySelector('#swap-langs');

const inTextArea = document.querySelector('#input-text');
const outTextArea = document.querySelector('#output-text');

const translateButton = document.querySelector('#submit-translation');

const infoStatus = document.querySelector('#info-status');
const infoError = document.querySelector('#info-error');

langSwapper.addEventListener('click', (e) => {
    const inLang = inLangSelect.value;
    const outLang = outLangSelect.value;

    inLangSelect.value = outLang;
    outLangSelect.value = inLang;
});

translateButton.addEventListener('click', async (e) => {
    const inLang = inLangSelect.value;
    const outLang = outLangSelect.value;
    const inText = inTextArea.value;

    if (inText === '') {
        outTextArea.value = '';
        return;
    }

    const payload = new FormData();
    payload.append('from_language', inLang);
    payload.append('to_language', outLang);
    payload.append('text', inText);

    hide(infoError);
    show(infoStatus);
    const resp = await fetch('/translate', {
        method: 'POST',
        body: payload
    });
    hide(infoStatus);

    if (resp.status === 200) {
        outTextArea.value = await resp.text();
    } else {
        infoError.textContent = 'Error: ' + await resp.text();
        show(infoError);
    }
});

function hide(elem) {
    elem.classList.add('hidden');
}

function show(elem) {
    elem.classList.remove('hidden');
}
