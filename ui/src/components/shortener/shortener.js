import React, { useState } from 'react'
import { Input } from '../input/input'
import { ShortInput } from '../shortInput/shortInput'
import { Button } from '../button/button';
import ShortApi from '../../api/shortAPI';

import './shortener.css'

export const Shortener = () => {
    const [longURI, setLongUri] = useState(null)
    const [short, setShort] = useState(null)

    const handleGenerateShort = (longURI, short) => {
        console.log(longURI, short);
        ShortApi
            .generateShort(longURI, short)
            .then(generatedShort => setShort(generatedShort))
    }

    return (
        <section>
            <p>
                <b>As a</b> smart ass
            </p>
            <p>
                <b>I want</b> <Input onChange={e => setLongUri(e.target.value)} className="long" placeholder="this long url" autoFocus/>
            </p>
            <p>
                <span>To be short </span>  
                <ShortInput domain="https://shortener.com/" className="short" onChange={e => setShort(e.target.value)} />
            </p>
            <p>
                <b>So that</b> it looks fancy. 
            </p>
            <Button text="Save" onClick={() => handleGenerateShort(longURI, short)}/>
        </section>)
}