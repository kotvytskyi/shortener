import React from 'react'
import { Input } from '../input/input'

export const ShortInput = ({ domain }) => (
    <span>
        <span>{ domain }</span> <Input />
    </span>
)