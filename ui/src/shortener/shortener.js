import React from 'react'
import { Input } from '../input/input'
import './shortener.css'

export const Shortener = () => (
    <section>
        <p>
            <b>As a</b> smart ass
        </p>
        <p>
            <b>I want</b> <Input className="long" placeholder="this long url" autoFocus/>
        </p>
        <p>
            To be short <span>https://shortener.com/ <Input className="short" /></span>
        </p>
        <p>
            <b>So that</b> it looks fancy. <button>Amen</button>
        </p>
    </section>
)