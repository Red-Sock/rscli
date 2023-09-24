import cls from './content.module.css'

import ReactMarkdown from 'react-markdown'
import {Footer} from "../footer/footer.tsx";
import {Sidebar} from "../sidebar/sidebar.tsx";


export function ContentWrapper  (props: { isOpen: any; }) {
    return (
        <div className={cls.ContentWrapper}>
            <div className={cls.Content}>
                <ReactMarkdown>
                    ## 2
                </ReactMarkdown>

                <Sidebar isOpen={props.isOpen}/>
            </div>

            <Footer/>
        </div>
    );
}
