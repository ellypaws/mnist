const apiUrl = '';

const canvas = document.getElementById('drawingCanvas');
const ctx = canvas.getContext('2d');
let drawing = false;
let strokes = false;

canvas.addEventListener('mousedown', startDrawing);
canvas.addEventListener('mouseup', stopDrawing);
canvas.addEventListener('mousemove', draw);
canvas.addEventListener('touchstart', startDrawing, {passive: false});
canvas.addEventListener('touchend', stopDrawing, {passive: false});
canvas.addEventListener('touchmove', draw, {passive: false});

const expectedInput = document.getElementById('expected');
const eraseButton = document.getElementById('eraseButton');
const resetButton = document.getElementById('resetButton');
const randomButton = document.getElementById('randomButton');
const additionButton = document.getElementById('additionButton');
const previewCheckbox = document.getElementById('previewCheckbox');
const additionDiv = document.getElementById('addition');
const augendElement = document.getElementById('augend');
const addendElement = document.getElementById('addend');
const sumElement = document.getElementById('sum');
const trainButton = document.getElementById('trainButton');
const submitButton = document.getElementById('submitButton');
const numberButtons = document.querySelectorAll('.number-button');
const predictionGraph = document.getElementById('predictionGraph').getContext('2d');
const previewCanvas = document.getElementById('previewCanvas');
const previewCtx = previewCanvas.getContext('2d');
let chart;

window.onload = randomizeExpected;

eraseButton.onclick = eraseCanvas
resetButton.onclick = resetCanvas;
// randomButton.onclick = randomizeExpected;
// additionButton.onclick = randomizeAddition;
previewCheckbox.onchange = drawTemplateNumber;
trainButton.onclick = startTraining;
submitButton.onclick = sendTrainingData;

numberButtons.forEach(button => button.onclick = () => {
    expectedInput.value = button.getAttribute('data-number');
    numberButtons.forEach(btn => btn.classList.remove('selected'));
    button.classList.add('selected');
    drawTemplateNumber()
});

function startDrawing(event) {
    drawing = true;
    draw(event);
}

function stopDrawing() {
    drawing = false;
    ctx.beginPath();
    if (strokes) {
        debounceSendDrawingToServer();
    }
    strokes = false;
}

function draw(event) {
    if (!drawing) return;

    const touch = event.touches ? event.touches[0] : event;
    strokes = true;
    ctx.lineWidth = 15;
    ctx.lineCap = 'round';
    ctx.strokeStyle = '#000';

    ctx.lineTo(touch.clientX - canvas.offsetLeft, touch.clientY - canvas.offsetTop);
    ctx.stroke();
    ctx.beginPath();
    ctx.moveTo(touch.clientX - canvas.offsetLeft, touch.clientY - canvas.offsetTop);
}

async function sendDrawingToServer() {
    const scaledCanvas = document.createElement('canvas');
    const scaledCtx = scaledCanvas.getContext('2d');
    scaledCanvas.width = 28;
    scaledCanvas.height = 28;

    scaledCtx.drawImage(canvas, 0, 0, 28, 28);

    // Process the image
    const processedCanvas = processImageForAPI(scaledCanvas);

    // Display the preview of the final output
    previewCtx.imageSmoothingEnabled = false;
    previewCtx.clearRect(0, 0, previewCanvas.width, previewCanvas.height);
    previewCtx.drawImage(processedCanvas, 0, 0, previewCanvas.width, previewCanvas.height);

    const imageData = processedCanvas.toDataURL('image/png');
    const expected = expectedInput.value ? parseInt(expectedInput.value) : null;

    const requestBody = {
        image: imageData,
        expected: expected
    };

    fetch(apiUrl + '/v1/predict', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestBody)
    })
        .then(response => response.json())
        .then(data => {
            updateButtonStates(data.prediction);
            displayPredictions(data);
            return data.correct || data.prediction === expected;
        })
        .catch(error => console.error('Error:', error));
}

function startTraining() {
    fetch(apiUrl + '/v1/train', {
        method: 'POST'
    })
        .then(response => response.blob())
        .then(blob => {
            console.log('Training finished successfully');
        })
        .catch(error => console.error('Error:', error));
}

async function sendTrainingData() {
    const scaledCanvas = document.createElement('canvas');
    const scaledCtx = scaledCanvas.getContext('2d');
    scaledCanvas.width = 28;
    scaledCanvas.height = 28;

    scaledCtx.drawImage(canvas, 0, 0, 28, 28);

    const processedCanvas = processImageForAPI(scaledCanvas);

    const imageData = processedCanvas.toDataURL('image/png');
    const expected = expectedInput.value ? parseInt(expectedInput.value) : null;

    if (!expected) {
        alert('Please select a number before submitting');
        return;
    }

    const correct = await sendDrawingToServer();

    const requestBody = {
        image: imageData,
        expected: expected,
        correct: correct,
    };

    fetch(apiUrl + '/v1/train', {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestBody)
    })
        .then(response => response.blob())
        .then(blob => {
            console.log('Training data sent successfully');
        })
        .catch(error => console.error('Error:', error));
}

function debounce(func, wait) {
    let timeout;
    return function (...args) {
        clearTimeout(timeout);
        timeout = setTimeout(() => func.apply(this, args), wait);
    };
}

const debounceSendDrawingToServer = debounce(await sendDrawingToServer, 500);

function displayPredictions(data) {
    const predictionsDiv = document.getElementById('predictions');
    const correctClass = data.correct ? 'correct-true' : 'correct-false';
    const correctSymbol = data.correct ? '✔️' : '❌';

    predictionsDiv.innerHTML = `
        <b>Prediction</b>: ${data.prediction}
        <b>Expected</b>: ${data.expected != null ? data.expected : 'N/A'}
        <p class="${correctClass}">Correct: ${correctSymbol}</p>
    `;

    const labels = Object.keys(data.predictions).map(key => `Digit ${key}`);
    const values = Object.values(data.predictions);

    if (chart) {
        chart.destroy();
    }

    chart = new Chart(predictionGraph, {
        type: 'bar',
        data: {
            labels: labels,
            datasets: [{
                label: 'Prediction Confidence',
                data: values,
                backgroundColor: 'rgba(75, 192, 192, 0.2)',
                borderColor: 'rgba(75, 192, 192, 1)',
                borderWidth: 1
            }]
        },
        options: {
            scales: {
                y: {
                    beginAtZero: true
                }
            }
        }
    });
}

function updateButtonStates(prediction) {
    numberButtons.forEach(button => {
        const buttonNumber = button.getAttribute('data-number');
        button.classList.remove('wrong', 'not-selected', 'correct');

        if (button.classList.contains('selected')) {
            if (prediction == buttonNumber) {
                button.classList.add('correct');
            } else {
                button.classList.add('wrong');
            }
        } else if (prediction == buttonNumber) {
            button.classList.add('not-selected');
        }
    });
}

function eraseCanvas() {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    previewCtx.clearRect(0, 0, previewCanvas.width, previewCanvas.height);
}

function resetCanvas() {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    previewCtx.clearRect(0, 0, previewCanvas.width, previewCanvas.height);
    expectedInput.value = '';
    numberButtons.forEach(button => {
        button.classList.remove('selected', 'wrong', 'not-selected', 'correct');
    });
    additionDiv.classList.add('hidden');
    augendElement.innerText = '';
    addendElement.innerText = '';
    sumElement.innerText = '';

    const predictionsDiv = document.getElementById('predictions');
    predictionsDiv.innerHTML = '';
    if (chart) {
        chart.destroy();
        chart = null;
    }
    randomizeExpected();
}

function randomizeExpected() {
    const randomValue = Math.floor(Math.random() * 10);
    expectedInput.value = randomValue;
    numberButtons.forEach(button => button.classList.remove('selected'));
    numberButtons[randomValue].classList.add('selected');
    drawTemplateNumber()
}

function randomizeAddition() {
    const randomAugend = Math.floor(Math.random() * 10);
    const randomAddend = Math.floor(Math.random() * (10 - randomAugend));
    expectedInput.value = randomAugend + randomAddend;
    numberButtons.forEach(button => button.classList.remove('selected'));
    numberButtons[randomAugend + randomAddend].classList.add('selected');
    augendElement.innerText = String(randomAugend);
    addendElement.innerText = String(randomAddend);
    sumElement.innerText = String(randomAugend + randomAddend);
    additionDiv.classList.remove('hidden');
}

function drawTemplateNumber() {
    const previewChecked = previewCheckbox.checked;
    if (!previewChecked) return;
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    if (strokes) {
        const savedImage = ctx.getImageData(0, 0, canvas.width, canvas.height);
        ctx.putImageData(savedImage, 0, 0);
    }

    ctx.font = "250px sans-serif";
    ctx.fillStyle = "rgba(200, 200, 200, 0.3)";
    ctx.textAlign = "center";
    ctx.textBaseline = "middle";
    ctx.fillText(expectedInput.value, canvas.width / 2, canvas.height / 2);
}

function removeTransparency(img) {
    const canvas = document.createElement('canvas');
    canvas.width = img.width;
    canvas.height = img.height;
    const ctx = canvas.getContext('2d');
    ctx.drawImage(img, 0, 0);

    const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
    const data = imageData.data;

    for (let i = 0; i < data.length; i += 4) {
        const alpha = data[i + 3] / 255;
        data[i] = data[i] * alpha + 255 * (1 - alpha);     // Red
        data[i + 1] = data[i + 1] * alpha + 255 * (1 - alpha); // Green
        data[i + 2] = data[i + 2] * alpha + 255 * (1 - alpha); // Blue
        data[i + 3] = 255; // Set alpha to 255 (no transparency)
    }

    ctx.putImageData(imageData, 0, 0);
    return canvas;
}

function invertColors(img) {
    const canvas = document.createElement('canvas');
    canvas.width = img.width;
    canvas.height = img.height;
    const ctx = canvas.getContext('2d');
    ctx.drawImage(img, 0, 0);

    const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
    const data = imageData.data;

    for (let i = 0; i < data.length; i += 4) {
        data[i] = 255 - data[i];     // Invert Red
        data[i + 1] = 255 - data[i + 1]; // Invert Green
        data[i + 2] = 255 - data[i + 2]; // Invert Blue
    }

    ctx.putImageData(imageData, 0, 0);
    return canvas;
}

function antiAlias(img) {
    const canvas = document.createElement('canvas');
    canvas.width = img.width;
    canvas.height = img.height;
    const ctx = canvas.getContext('2d');
    ctx.drawImage(img, 0, 0);

    const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
    const data = imageData.data;
    const width = canvas.width;
    const height = canvas.height;

    // Kernel for a simple blur filter
    const kernel = [
        1 / 16, 2 / 16, 1 / 16,
        2 / 16, 4 / 16, 2 / 16,
        1 / 16, 2 / 16, 1 / 16
    ];
    const half = Math.floor(Math.sqrt(kernel.length) / 2);

    const applyKernel = (x, y) => {
        let r = 0, g = 0, b = 0, a = 0;
        for (let ky = -half; ky <= half; ky++) {
            for (let kx = -half; kx <= half; kx++) {
                const posX = Math.min(width - 1, Math.max(0, x + kx));
                const posY = Math.min(height - 1, Math.max(0, y + ky));
                const offset = (posY * width + posX) * 4;
                const weight = kernel[(ky + half) * 3 + (kx + half)];
                r += data[offset] * weight;
                g += data[offset + 1] * weight;
                b += data[offset + 2] * weight;
                a += data[offset + 3] * weight;
            }
        }
        return [r, g, b, a];
    };

    const outputData = new Uint8ClampedArray(data.length);
    for (let y = 0; y < height; y++) {
        for (let x = 0; x < width; x++) {
            const [r, g, b, a] = applyKernel(x, y);
            const offset = (y * width + x) * 4;
            outputData[offset] = r;
            outputData[offset + 1] = g;
            outputData[offset + 2] = b;
            outputData[offset + 3] = a;
        }
    }

    ctx.putImageData(new ImageData(outputData, width, height), 0, 0);
    return canvas;
}

function contrast(img, strength) {
    const canvas = document.createElement('canvas');
    canvas.width = img.width;
    canvas.height = img.height;
    const ctx = canvas.getContext('2d');
    ctx.drawImage(img, 0, 0);

    const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
    const data = imageData.data;

    const factor = (259 * (strength + 255)) / (255 * (259 - strength));

    for (let i = 0; i < data.length; i += 4) {
        data[i] = Math.min(255, factor * (data[i] - 128) + 128);     // Red
        data[i + 1] = Math.min(255, factor * (data[i + 1] - 128) + 128); // Green
        data[i + 2] = Math.min(255, factor * (data[i + 2] - 128) + 128); // Blue
    }

    ctx.putImageData(imageData, 0, 0);
    return canvas;
}

processors = [
    removeTransparency,
    invertColors,
    // antiAlias,
    img => contrast(img, 0.5)
];

function processImageForAPI(img) {
    for (const processor of processors) {
        img = processor(img);
    }
    return img;
}
