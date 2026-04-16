const getDeepWebGL = () => {
    const canvas = document.createElement('canvas');
    const gl = canvas.getContext('webgl') || canvas.getContext('experimental-webgl');
    if (!gl) return {};

    const debugInfo = gl.getExtension('WEBGL_debug_renderer_info');
    return {
        // vendor and model
        renderer: debugInfo ? gl.getParameter(debugInfo.UNMASKED_RENDERER_WEBGL) : 'unknown',
        vendor: debugInfo ? gl.getParameter(debugInfo.UNMASKED_VENDOR_WEBGL) : 'unknown',

        // hardware limits
        maxTextureSize: gl.getParameter(gl.MAX_TEXTURE_SIZE), // texture size (ex. 16384)
        maxRenderBufferSize: gl.getParameter(gl.MAX_RENDERBUFFER_SIZE),
        maxAttributes: gl.getParameter(gl.MAX_VERTEX_ATTRIBS),

        // float precision
        shaderPrecision: gl.getShaderPrecisionFormat(gl.FRAGMENT_SHADER, gl.HIGH_FLOAT).precision,

        // supported extensions
        extensions: gl.getSupportedExtensions().length
    };
};

const getCanvasFingerprint = () => {
    const canvas = document.createElement('canvas');
    const ctx = canvas.getContext('2d');
    canvas.width = 240;
    canvas.height = 60;

    ctx.textBaseline = "top";
    ctx.font = "14px 'Arial'";
    ctx.fillStyle = "#f60";
    ctx.fillRect(125, 1, 62, 20);
    ctx.fillStyle = "#069";
    ctx.fillText("Unique-ID <canvas> 😃", 2, 15);
    ctx.fillStyle = "rgba(102, 204, 0, 0.7)";
    ctx.fillText("Unique-ID <canvas> 😃", 4, 17);

    return canvas.toDataURL();
};

const getAudioFingerprint = async () => {
    const context = new (window.OfflineAudioContext || window.webkitOfflineAudioContext)(1, 44100, 44100);
    const oscillator = context.createOscillator();
    oscillator.type = 'triangle';
    oscillator.frequency.setValueAtTime(10000, context.currentTime);

    const compressor = context.createDynamicsCompressor();
    compressor.threshold.setValueAtTime(-50, context.currentTime);
    compressor.knee.setValueAtTime(40, context.currentTime);
    compressor.ratio.setValueAtTime(12, context.currentTime);

    oscillator.connect(compressor);
    compressor.connect(context.destination);
    oscillator.start(0);

    const renderedBuffer = await context.startRendering();
    const samples = renderedBuffer.getChannelData(0).slice(4500, 5000); // piece of data
    let sum = 0;
    for (let i = 0; i < samples.length; i++) {
        sum += Math.abs(samples[i]);
    }
    return sum;
};

function cyrb53(str, seed = 0) {
    let h1 = 0xdeadbeef ^ seed, h2 = 0x41c6ce57 ^ seed;
    for (let i = 0, ch; i < str.length; i++) {
        ch = str.charCodeAt(i);
        h1 = Math.imul(h1 ^ ch, 2654435761);
        h2 = Math.imul(h2 ^ ch, 1597334677);
    }
    h1 = Math.imul(h1 ^ (h1 >>> 16), 2246822507) ^ Math.imul(h2 ^ (h2 >>> 13), 3266489909);
    h2 = Math.imul(h2 ^ (h2 >>> 16), 2246822507) ^ Math.imul(h1 ^ (h1 >>> 13), 3266489909);
    return 4294967296 * (2097151 & h2) + (h1 >>> 0);
}

const set = (id, val) =>
    document.getElementById(id).textContent = JSON.stringify(val)

export async function getVisitorId() {
    const webgl = getDeepWebGL()
    set("fp-webgl", webgl)

    const canvas = getCanvasFingerprint()
    set("fp-canvas", canvas.slice(0, 30) + "...") // dataURL длинный

    const audio = await getAudioFingerprint()
    set("fp-audio", audio)

    const cores = navigator.hardwareConcurrency || 0
    set("fp-cores", cores)

    const ram = navigator.deviceMemory || 0
    set("fp-ram", ram)

    const screen_ = `${screen.width}x${screen.height}`
    set("fp-screen", screen_)

    const id = cyrb53(JSON.stringify({ webgl, canvas, audio, cores, ram, screen: screen_ }))
    const visitorId = id.toString(36)
    set("fp-result", visitorId)

    return visitorId
}

// export async function getVisitorId() {
//     const components = {
//         webgl: getDeepWebGL(),
//         canvas: getCanvasFingerprint(),
//         audio: await getAudioFingerprint(),
//         cores: navigator.hardwareConcurrency || 0,
//         ram: navigator.deviceMemory || 0,
//         screen: `${screen.width}x${screen.height}`
//     };

//     const id = cyrb53(JSON.stringify(components));
//     return id.toString(36);
// }

; (async () => {
    const visitorId = await getVisitorId()

    // const response = await fetch("/api/check-vote", {
    //     method: "POST",
    //     headers: { "Content-Type": "application/json" },
    //     body: JSON.stringify({ visitorId })
    // })

    // const { isVotedToday } = await response.json()
    // console.log(isVotedToday)
})()
