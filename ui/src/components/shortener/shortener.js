import React, { useState, useEffect } from 'react'
import { Input } from '../input/input'
import { ShortInput } from '../shortInput/shortInput'
import { ShortLink } from '../shortLink/shortLink'
import { Button } from '../button/button';
import ShortApi from '../../api/shortAPI';

import './shortener.css'

export const Shortener = () => {
    const [longURI, setLongUri] = useState(null)
    const [short, setShort] = useState(null)

    const handleGenerateShort = (originalURI, short) => {
        ShortApi
            .generateShort(originalURI, short)
            .then(generatedShort => setShort(generatedShort))
    }

    return (
        <section>
            <p>
                <b>As a</b> smart ass
            </p>
            <p>
                <b>I want</b> <Input onChange={e => console.log(e.target.val)} className="long" placeholder="this long url" autoFocus/>
            </p>
            <p>
                <span>To be short </span>  
                { !short ?  
                    <ShortInput domain="https://shortener.com/" className="short" /> :
                    <ShortLink link={ short }></ShortLink>
                }
            </p>
            <p>
                <b>So that</b> it looks fancy. 
            </p>
            <Button text="Save" onClick={() => generateShort("https://shortener.com/test")}/>
        </section>)
}