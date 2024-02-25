import { useNavigate } from "react-router-dom";
import { useDispatch } from "react-redux";
import './styles.css';
import config from '../config.json';
import Axios from "axios";
import { GET_USER_ID } from "../action_types";

// Set action type
//const GET_NEW_USER_ID = "set_user_id";

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
            dispatch({ type: GET_USER_ID, payload: newUserID });

            navigate('/chat');
    
        }
        catch (error) {
            console.error("Error fetching data:", error);
        }
    };

  return (
  <div>
    <button className='btn btn-primary btn-sx' onClick={handleClick}>
      {children}
    </button>
  </div>
  )
}

export default LoginButton