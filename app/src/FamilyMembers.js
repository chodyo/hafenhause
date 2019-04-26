import React from 'react';
import axios from 'axios';
import ColorfulDigitCountdown from './ColorfulDigitCountdown';

const ulStyle = {
    listStyle: "none"
}

const liStyle = {
    padding: "1em",
    float: "left"
}

const dateDisplayOpts = { month: "short", day: "numeric", hour: "numeric", minute: "numeric", second: "numeric", millisecond: "numeric" }


export default class FamilyMembers extends React.Component {

    state = {
        members: []
    }

    componentDidMount() {
        axios.get(`https://us-central1-hafenhause.cloudfunctions.net/Bedtime`)
            .then(res => {
                const members = res.data;
                this.setState({ members });
            });
    }

    render() {
        return (
            <ul style={ulStyle}>
                {this.state.members.map(member => {
                    let updated = new Date(member.updated);
                    let bedtime = new Date(updated);

                    bedtime.setHours(member.hour, member.minute, 0, 0);

                    // if a time gets updated after the bedtime has passed, it
                    // was referring to the next day
                    if (bedtime < updated) {
                        bedtime.setDate(bedtime.getDate() + 1);
                    }

                    return (
                        <li key={member.name} style={liStyle}>
                            <div>{member.name}</div>
                            <div>{bedtime.toLocaleString("en-US", dateDisplayOpts)}</div>
                            <ColorfulDigitCountdown date={bedtime} />
                        </li>
                    );
                })}
            </ul>
        );
    }
}
