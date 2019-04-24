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

                    bedtime.setHours(member.hour);
                    bedtime.setMinutes(member.minute);

                    // if a time gets updated after the bedtime has passed, it
                    // was referring to the next day
                    if (bedtime < updated) {
                        bedtime.setDate(bedtime.getDate() + 1);
                    }

                    return <li style={liStyle}>
                        <div>{member.name}</div>
                        <div>{bedtime.toLocaleString("en-US", dateDisplayOpts)}</div>
                        <Countdown
                            date={bedtime}
                            daysInHours={true}
                        />
                    </li>
                })}
            </ul>
        );
    }
}
