import './index.module.css'

import React from 'react'
import ReactDOM from 'react-dom/client'

import {Home} from './pages/home/home';
import {BrowserRouter} from "react-router-dom";

ReactDOM.createRoot(document.getElementById('root')!).render(
    <React.StrictMode>
        <link href="https://fonts.googleapis.com/icon?family=Comfortaa" rel="stylesheet"/>
        <BrowserRouter>
            <Home/>
        </BrowserRouter>
    </React.StrictMode>
)
