import {useState} from 'react';
import logo from './assets/images/logo.png';
import './App.css';
import {Export} from "../wailsjs/go/main/App";

function App() {
    const [name, setName] = useState('');
    const [resultText, setResultText] = useState(undefined);
    const [loading, setLoading] = useState(false);

    const exportTxs = async () => {
        setLoading(true);
        setResultText(undefined);
        const result = await Export(name, "--logs --articulate");
        setResultText(result);
        setLoading(false);
    }

    return (
        <div id="App">
            <img src={logo} alt="logo"/>
            <div className="address-prompt">
                Enter an address or ENS name below 👇
            </div>
            <div className="input-box">
                <input 
                    className="input" 
                    onChange={(e) => setName(e.target.value)} 
                    onKeyDown={(e) => e.key === 'Enter' && exportTxs()}
                    value={name}
                    placeholder="trueblocks.eth" 
                    autoComplete="off" 
                    name="input" 
                    autoFocus
                />
                <button 
                    className="btn" 
                    onClick={exportTxs}
                    disabled={loading || name === ''}
                >
                    Export
                </button>
            </div>
            {loading &&
                <div className="result">
                    Loading...
                </div>
            }
            {resultText !== undefined &&
                <div className='result'>
                    {resultText}
                </div>
            }
        </div>
    )
}

export default App
