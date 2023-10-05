import cls from './content.module.css'

import ReactMarkdown from 'react-markdown'
import React from "react";

export const ContentWrapper = (props: {content: string}) => {
    return (
        <div className={cls.ContentWrapper}>
            <ReactMarkdown
                children={props.content}/>
        </div>
    )
}

