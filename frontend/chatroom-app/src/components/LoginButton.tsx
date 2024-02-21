import { useNavigate } from "react-router-dom";
import { useDispatch } from "react-redux";
import '../styles.css';
import config from '../../config.json';
import Axios from "axios";
import store from "../store.tsx";

// Set action type
const GET_NEW_USER_ID = "get/newuserID";

interface Props {
    children: string;
    // '?' gives color a value is optional, similar to default param in c++
    // this allows for color to only represent either primary secondary or danger
    //color?: 'primary' | 'secondary' | 'danger';
    //color?: string;
    //onClick: () => void;
}

const LoginButton = ({ children }: Props) => {

    const navigate = useNavigate();
    const dispatch = useDispatch();

    const handleClick = async () => {
        try {
            // to get userID
            const serverHost = config.server.host;
            const serverPort = config.server.port;
            const path = "/newuser";
            const url = `http://${serverHost}:${serverPort}${path}`;
            
            // PUT request to server
            const response = await Axios.put(url);
            const newUserID = response.data;

            // Dispatch an action to update the store with new userID
            dispatch({ type: GET_NEW_USER_ID, payload: newUserID });

            navigate('/Chat');
    
        }
        catch (error) {
            console.error("Error fetching data:", error);
        }
    };

    return (
  <div style={{ position: 'relative', minHeight: '100vh', display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center' }}>
    <div style={{ marginBottom: '20px' }}>
      <img src="../../public/chatroomLogo.gif" alt="Chatroom Logo" style={{ maxWidth: '1000%', maxHeight: '50vh', position: 'absolute', top: '40%', left: '-100%', right: '50%', transform: 'translate(-50%, -50%)' }} />
    </div>
    <div>
      <button className='btn btn-primary btn-sx' onClick={handleClick}>
        {children}
      </button>
    </div>
  </div>
  )
}

export default LoginButton