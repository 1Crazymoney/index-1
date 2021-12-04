import {Host} from "../../components/config"

export default function handler(req, res) {
    return new Promise((resolve, reject) => {
        fetch(Host + "/hello").then(res => res.json()).then(data => {
            res.status(200).json(data)
            res.end()
            resolve()
        }).catch(error => {
            res.status(500).json(error)
            res.end()
            reject(error)
        })
    })
}
