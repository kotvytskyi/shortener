const ShortApi = {};

ShortApi.generateShort = (originalURI, short) => {
    return fetch('/api/shorts',
    {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ url: originalURI, short: short })
    }).then(response => response.json());
};

ShortApi.isShortExist = (short) => {
    return fetch(`/api/shorts/${short}`).then(response => response.status == 200);
};

export default ShortApi