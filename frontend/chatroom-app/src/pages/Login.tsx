import React, { useEffect } from 'react';
import LoginButton from '.././components/LoginButton';
import config from '../../config.json'
import Axios from "axios";
import { useState } from "react";

export interface LoginProps { }

const LoginPage: React.FunctionComponent<LoginProps> = (props) => {

    // make state to store the user
    //const [userID, setuserID] = useState("");

    const serverHost = config.server.host;
    const serverPort = config.server.port;
    const path = "/newuser";
    const url = `http://${serverHost}:${serverPort}${path}`;

    //useEffect(() => {
    Axios.put(url).then((respond) => {
            console.log(url)
            console.log(respond.data);
            //setuserID(respond.data);
        })
            .catch((error) => {
                console.error("Error fetching data:", error);
        });
    //}, []);

    return (
        <div className="LoginPage">
            <LoginButton>Log In</LoginButton>
        </div>
    );
};

export default LoginPage;