import React from 'react'
import { Input } from '../input/input'
import './shortInput.css'

export const ShortInput = ({ domain, ...props }) => (
    <span>
        <span>{ domain }</span> 
        <div className="short-row">
            <Input {...props} />
        </div>
    </span>
)