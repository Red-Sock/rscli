import { motion } from 'framer-motion';
import cls from './sidebar.module.css'
export function Sidebar (props: { isOpen: any; }) {

    return (
        <motion.div
            className={cls.Sidebar}
            layout
            data-isOpen={props.isOpen}
        >
            <div className={cls.Content}>
                <ul>
                    <li>123faefavtrsvtreagrvaerfgerg</li>
                    <li>1234</li>
                </ul>
            </div>
        </motion.div>
    );
}
