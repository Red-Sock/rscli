import cls from './content.module.css'

import ReactMarkdown from 'react-markdown'
import {memo, useState} from "react";


export const ContentWrapper = memo(() => {
        const contentURL = "https://raw.githubusercontent.com/Red-Sock/rscli/docs/docs/";
        const projectPath = "rscli";
        const docsFolder = "docs"
        const fileExtension = ".md"

        const resourcePath = window.location.href.substring(
            window.location.href.indexOf(projectPath)+projectPath.length+1
        )


        const [content, setContent] = useState("")
        const resourceURL = [contentURL, resourcePath+fileExtension].join("/")

        fetch(resourceURL).then(async (response) => {
            if (response.ok) {
                setContent(await response.text())
            } else {
                window.location.replace([window.location.href.substring(0, window.location.href.indexOf(projectPath)-1), projectPath, "home"].join("/"));
            }
        })

        return (
            <div className={cls.ContentWrapper}>
                <div className={cls.ContentField}>
                    <div className={cls.Content}>
                        <ReactMarkdown
                            children={content}/>
                    </div>
                </div>
            </div>
        );
    }
)
