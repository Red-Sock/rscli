import cls from './home.module.css';
import {Header} from "../../sections/header/header.tsx";
import {Footer} from "../../sections/footer/footer.tsx";
import {useState} from "react";
import {Sidebar} from "../../sections/header/sidebar/sidebar.tsx";
import {ContentDisplay} from "../../sections/content/content.tsx";
export function Home() {
  const [isSideMenuOpen, setIsSideMenuOpen] = useState(false)
    return (
    <div className={cls.Home}>
        <Header
            isOpen={isSideMenuOpen}
            setIsOpen={setIsSideMenuOpen}
        />
        <div className={cls.contentWrap}>
            <Sidebar
             isOpen={isSideMenuOpen}
            />
            <ContentDisplay/>
        </div>
        <Footer/>
    </div>
  )
}
