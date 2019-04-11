import React from 'react';
import axios from 'axios';
import Countdown from 'react-countdown-now';

const ulStyle = {
    listStyle: "none"
}

const liStyle = {
    padding: "1em",
    float: "left"
}

export default class FamilyMembers extends React.Component {
    state = {
        members: []
    }

    componentDidMount() {
        axios.get(`https://us-central1-hafenhause.cloudfunctions.net/GetBedtimes`)
            .then(res => {
                const members = res.data;
                this.setState({ members });
            });
    }

    render() {
        return (
            <ul style={ulStyle}>
                {this.state.members.map(member => {
                    let options = { month: "short", day: "numeric", hour: "numeric", minute: "numeric", second: "numeric", millisecond: "numeric" }
                    let date = new Date(member.date).toLocaleString("en-US", options);
                    return <li style={liStyle}>
                        <div>{member.name}</div>
                        <div>{date}</div>
                        <Countdown
                            date={member.date}
                            daysInHours={true}
                        />
                    </li>
                })}
            </ul>
        );
    }
}