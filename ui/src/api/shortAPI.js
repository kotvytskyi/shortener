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

export default ShortApi