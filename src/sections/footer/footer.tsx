import cls from './footer.module.css';

import {useState} from "react";

import {motion} from 'framer-motion';

export function Footer() {
    const [isOpen, setIsOpen] = useState(false)

    return (
        <motion.div
            layout
            className={cls.Footer}
            data-isOpen={isOpen}
        >
            <div className={cls.ButtonContainer}>
                <div
                    onClick={() => setIsOpen(!isOpen)}
                    className={cls.OpenButton}/>
            </div>
        </motion.div>
    );
}

