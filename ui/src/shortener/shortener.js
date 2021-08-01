import React, { useState, useEffect } from 'react'
import { Input } from '../input/input'
import { ShortInput } from '../shortInput/shortInput'
import { ShortLink } from '../shortLink/shortLink'
import { Button } from '../button/button';

import './shortener.css'

export const Shortener = () => {
    const [short, generateShort] = useState(null)

    return (
        <section>
            <p>
                <b>As a</b> smart ass
            </p>
            <p>
                <b>I want</b> <Input className="long" placeholder="this long url" autoFocus/>
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