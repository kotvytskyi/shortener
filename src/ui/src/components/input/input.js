import React from 'react'
import './input.css'

export const Input = ({ isInvalid, errorMessage, ...props }) => (
    <div className="container">
        <input {...props} />
        {isInvalid ?
            <div className="error_message">
                <span>{errorMessage}</span>
            </div> : null}
    </div>
)