import React, { useState } from 'react'
import { Input } from '../input/input'
import { ShortInput } from '../shortInput/shortInput'
import { Button } from '../button/button';
import ShortApi from '../../api/shortAPI';
import validator from 'validator';

import './shortener.css'

export const Shortener = () => {
    const [uri, setUri] = useState(null);
    const [uriValid, setUriValid] = useState(false)

    const [uriErrorMessage, setUriErrorMessage] = useState(null)
    const [shortErrorMessage, setShortErrorMessage] = useState(null)

    const [short, setShort] = useState(null)
    const [shortValid, setShortValid] = useState(false)

    const handleUriChange= (input) => {
        setUri(input.target.value);
    }

    const handleGenerateShort = (longURI, short) => {
        ShortApi
            .generateShort(longURI, short)
            .then(generatedShort => setShort(generatedShort))
    }

    const validateUri = uri => ({
        valid: uri && validator.isURL(uri),
        message: "URL is invalid"
    })
    const validateShort = short => ({
        valid: short && validator.isLength(short, { min: 3, max: 10 }),
        message: "Short length must be between 3 and 10 characters"
    })

    const handleSubmit = (_) => {
        const uriValidation = validateUri(uri)
        const shortValidation = validateShort(short)

        setUriValid(uriValidation.valid)
        setUriErrorMessage(uriValidation.message)

        setShortValid(shortValidation.valid)
        setShortErrorMessage(shortValidation.message)

        if (uriValidation.valid && shortValidation.valid) {
            handleGenerateShort(uri, short)
        }
    }

    return (
        <section>
            <p>
                <b>As a</b> smart ass
            </p>
            <div className="row">
                <b>I want </b> 
                <Input 
                    isInvalid={!uriValid}
                    errorMessage={uriErrorMessage}
                    onChange={handleUriChange}
                    className="long"
                    placeholder="this long url"
                    autoFocus/>
            </div>
            <div className="row">
                <span>To be short </span>  
                <ShortInput
                    domain="https://shortener.com/"
                    className="short"
                    isInvalid={!shortValid}
                    errorMessage={shortErrorMessage}
                    onChange={e => setShort(e.target.value)} />
            </div>
            <p>
                <b>So that</b> it looks fancy. 
            </p>
            <Button text="Save" onClick={handleSubmit}/>
        </section>)
}