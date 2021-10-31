import React from 'react'
import './button.css'

export const Button = ({ text, ...props }) => (
    <button {...props}>{text}</button>
)