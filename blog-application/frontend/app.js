import Navigo from 'navigo'
import { AuthServiceClient,LoginRequest,AuthUserRequest,SignupRequest } from './proto/services_grpc_web_pb'

const router = new Navigo()
const authClient = new AuthServiceClient("http://localhost:9001")

router
.on("/",function(){
    document.body.innerHTML = "Home"
})
.on("/login",function(){
    document.body.innerHTML = ""
    const loginDiv = document.createElement('div')
    loginDiv.classList.add("login-div")

    const loginLabel = document.createElement("h1")
    loginLabel.innerText = "Login"
    loginDiv.appendChild(loginLabel)

    const loginForm = document.createElement("form")
    
    const loginInput =  document.createElement("input")
    loginInput.setAttribute('type','text')
    loginInput.setAttribute('placeholder','Enter username')
    loginForm.append(loginInput)
    
    const passwordInput =  document.createElement("input")
    passwordInput.setAttribute('type','password')
    passwordInput.setAttribute('placeholder','Enter password')
    loginForm.append(passwordInput)

    loginForm.append(document.createElement("div"))

    const submitButton =  document.createElement("button")
    submitButton.innerText = "Login"
    loginForm.append(submitButton)

    loginForm.addEventListener('submit', event => {
        event.preventDefault()
        let req = new LoginRequest()
        req.setLogin(loginInput.value)
        req.setPassword(passwordInput.value)
        authClient.login(req,{},(err,res) => {
            if(err){
                return alert(err.message)
            }
            //console.log(res.getToken())
            localStorage.setItem('token',res.getToken())
            req = new AuthUserRequest()
            req.setToken(res.getToken())
            authClient.authUser(req,{},(err,res)=>{
                if(err){
                    return alert(err.message)
                }
                const user = {id:res.getId(),username:res.getUsername(),email:res.getEmail()}
                localStorage.setItem('user',JSON.stringify(user))
            })
        })
        var activeUser = JSON.parse(localStorage.getItem('user'))
        alert(activeUser.username+" is logged in!")
    })

    loginDiv.appendChild(loginForm)
    
    document.body.appendChild(loginDiv)
})
.on("/signup",function(){

    document.body.innerHTML = ""
    const signupDiv = document.createElement('div')
    signupDiv.classList.add("signup-div") 

    const signupLabel = document.createElement("h1")
    signupLabel.innerText = "Signup"
    signupDiv.appendChild(signupLabel)

    const signupForm = document.createElement("form")
    
    const signupInput =  document.createElement("input")
    signupInput.setAttribute('type','text')
    signupInput.setAttribute('placeholder','Enter username')
    signupForm.append(signupInput)
    
    const passwordInput =  document.createElement("input")
    passwordInput.setAttribute('type','password')
    passwordInput.setAttribute('placeholder','Enter password')
    signupForm.append(passwordInput)

    const emailInput =  document.createElement("input")
    emailInput.setAttribute('type','email')
    emailInput.setAttribute('placeholder','Enter email')
    signupForm.append(emailInput)

    signupForm.append(document.createElement("div"))

    const submitButton =  document.createElement("button")
    submitButton.innerText = "Signup"
    signupForm.append(submitButton)

    signupForm.addEventListener('submit', event => {
        event.preventDefault()
        let req = new SignupRequest()
        req.setUsername(signupInput.value)
        req.setPassword(passwordInput.value)
        req.setEmail(emailInput.value)
        authClient.signup(req,{},(err,res) => {
            if(err){
                return alert(err.message)
            }
            //console.log(res.getToken())
            localStorage.setItem('token',res.getToken())
            req = new AuthUserRequest()
            req.setToken(res.getToken())
            authClient.authUser(req,{},(err,res)=>{
                if(err){
                    return alert(err.message)
                }
                const user = {id:res.getId(),username:res.getUsername(),email:res.getEmail()}
                localStorage.setItem('user',JSON.stringify(user))
            })
        })
        var activeUser = JSON.parse(localStorage.getItem('user'))
        alert("Registered!!")
    })

    signupDiv.appendChild(signupForm)

    document.body.appendChild(signupDiv)
})
.resolve()