import React from 'react'
import './shortLink.css'

export const ShortLink = ({ link, ...props }) => (
    <a href={link} {...props}>{link}</a>
)