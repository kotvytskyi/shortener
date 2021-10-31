const baseURI = 'http://localhost:80/api'

const ShortApi = {};

ShortApi.generateShort = (originalURI, short) => {
    return fetch(baseURI + '/shorts',
    {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ url: originalURI, short: short })
    }).then(response => response.json());
};

export default ShortApi