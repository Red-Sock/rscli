import cls from './content.module.css'

import ReactMarkdown from 'react-markdown'
import {useState} from "react";

import {doStaff} from '../../services/file-fetcher';


export const ContentWrapper = () => {

    const path = doStaff();
    console.log('res url', path.resourceURL);
    console.log('ref url', path.refreshURl);

    const [content, setContent] = useState("# Content is loading...")

    if (path.resourceURL.length === 0) {
        window.location.replace(path.refreshURl);
    }

    fetch(path.resourceURL).then(async (response) => {
        if (response.ok) {
            setContent(await response.text())
        }
        window.location.replace(path.refreshURl);
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

