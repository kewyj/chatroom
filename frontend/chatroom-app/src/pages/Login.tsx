import React from 'react';
import LoginButton from '.././components/LoginButton';

export interface LoginProps {}

const LoginPage: React.FunctionComponent<LoginProps> = (props) => {
    return (
        <div>
            <LoginButton>Log In</LoginButton>
        </div>
    );
};

export default LoginPage;