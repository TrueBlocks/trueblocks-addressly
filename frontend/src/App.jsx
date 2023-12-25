import {useState} from 'react';
import logo from './assets/images/logo.png';
import './App.css';
import {Export} from "../wailsjs/go/main/App";

function App() {
    const [resultText, setResultText] = useState("Please enter your name below 👇");
    const [name, setName] = useState('');
    const updateName = (e) => setName(e.target.value);
    const updateResultText = (result) => setResultText(result);

    function exportTxs() {
        Export(name).then(updateResultText);
    }

    return (
        <div id="App">
            <img src={logo} id="logo" alt="logo"/>
            <div id="result" className="result"><pre>{resultText}</pre></div>
            <div id="input" className="input-box">
                <input id="name" className="input" onChange={updateName} placeholder="trueblocks.eth" autoComplete="off" name="input" type="text"/>
                <button className="btn" onClick={exportTxs}>Export</button>
            </div>
        </div>
    )
}

export default App
