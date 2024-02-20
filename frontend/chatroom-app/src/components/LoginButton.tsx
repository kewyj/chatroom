import { useNavigate } from "react-router-dom";

interface Props {
    children: string;
    // '?' gives color a value is optional, similar to default param in c++
    // this allows for color to only represent either primary secondary or danger
    //color?: 'primary' | 'secondary' | 'danger';
    //color?: string;
    onClick: () => void;
}

const LoginButton = ({ children }: Props) => {
    const navigate = useNavigate();

    return (
      <button className='btn btn-primary'
      onClick={() => navigate('/Chat')}>
          {children}
      </button>
  )
}

export default LoginButton