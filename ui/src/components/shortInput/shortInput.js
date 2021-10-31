import React from 'react'
import { Input } from '../input/input'

export const ShortInput = ({ domain, ...props }) => (
    <span>
        <span>{ domain }</span> <Input {...props} />
    </span>
)