import { useNavigate } from "react-router-dom";
import '../styles.css'

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

    return (
        <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '100vh' }}>
            <img src="../../public/chatroomLogo.gif" alt="Chatroom Logo" style={{ maxWidth: '100%', maxHeight: '50vh', position: 'absolute', top: '35%', transform: 'translateX(2%)' }}/>
            <button className='btn btn-primary btn-sx'
                onClick={() => navigate('/Chat')}>
                {children}
            </button>
        </div>
  )
}

export default LoginButton