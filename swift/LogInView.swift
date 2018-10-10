
import UIKit

enum LogInResult {
    case emailFailure
    case passwordFailure
    case generalFailure
    case success
}

class LogInView: UIView {
    @IBOutlet weak var emailLabel: UILabel!
    @IBOutlet weak var emailField: UITextField!
    @IBOutlet weak var passwordLabel: UILabel!
    @IBOutlet weak var passwordField: UITextField!
    @IBOutlet weak var loginButton: UIButton!
    @IBOutlet weak var forgotButton: UIButton!
    @IBOutlet weak var signUpButton: UIButton!
    @IBOutlet weak var buttonSpinner: UIActivityIndicatorView!
    weak var currentField: UITextField?

    var loginCallback: ((_ email: String, _ password: String, _ completion: @escaping (LogInResult) -> Void) -> Void)?
    var showForgotPasswordCallback: (() -> Void)?
    var showSignUpCallback: (() -> Void)?

    override func awakeFromNib() {
        super.awakeFromNib()

        emailLabel.attributedText = NSAttributedString(string: NSLocalizedString("log-in.label.email", comment: "Sign In"), attributes: OnboardingViewController.fieldAttributes)
        passwordLabel.attributedText = NSAttributedString(string: NSLocalizedString("log-in.label.password", comment: "Sign In"), attributes: OnboardingViewController.fieldAttributes)
        loginButton.setTitle(NSLocalizedString("log-in.button.login", comment: "Sign In"), for: .normal)
        forgotButton.setTitle(NSLocalizedString("log-in.button.forgot", comment: "Sign In"), for: .normal)
        signUpButton.setTitle(NSLocalizedString("log-in.button.signup", comment: "Sign In"), for: .normal)
    }

    // MARK: - Actions
    
    @IBAction func loginButtonTapped(_ sender: UIButton) {
        guard let login = loginCallback else { return }
        guard let email = emailField.text else { return }
        guard let password = passwordField.text else { return }

        setLoginButtonEnabled(false)
        login(email, password) { result in
            switch result {
            case .emailFailure:
                self.emailLabel.textColor = Style.Color.red
                self.passwordLabel.textColor = Style.Color.black
            case .passwordFailure:
                self.emailLabel.textColor = Style.Color.black
                self.passwordLabel.textColor = Style.Color.red
            case .generalFailure:
                self.emailLabel.textColor = Style.Color.black
                self.passwordLabel.textColor = Style.Color.black
            default:
                break
            }
            self.setLoginButtonEnabled(true)
        }
    }
    
    @IBAction func showForgotPassword(_ sender: UIButton) {
        self.resignFirstResponder()
        
        guard let cb = showForgotPasswordCallback else { return }
        cb()
    }
    
    @IBAction func showSignUp(_ sender: UIButton) {
        self.resignFirstResponder()
        
        guard let cb = showSignUpCallback else { return }
        cb()
    }
    
    // MARK: - Private helpers
    
    private func setLoginButtonEnabled(_ enabled: Bool) {
        loginButton.isEnabled = enabled
        loginButton.backgroundColor = enabled ? Style.Color.black : Style.Color.gray
        if enabled {
            buttonSpinner.stopAnimating()
        } else {
            buttonSpinner.startAnimating()
        }
    }
}

extension LogInView: MaximumSizeable {
    var maxSize: CGSize {
        get {
            return CGSize(width: 0, height: 352)
        }
    }
}

extension LogInView: UITextFieldDelegate {
    func textFieldDidBeginEditing(_ textField: UITextField) {
        updateLabelColorsForTextField(textField)
        currentField = textField
    }
    
    func textField(_ textField: UITextField, shouldChangeCharactersIn range: NSRange, replacementString string: String) -> Bool {
        updateLabelColorsForTextField(textField)
        return true
    }

    func textFieldShouldReturn(_ textField: UITextField) -> Bool {
        let text = textField.text ?? ""
        
        if textField == emailField {
            if text.isValidEmail {
                emailLabel.textColor = Style.Color.black
                passwordField.becomeFirstResponder()
                return true
            }
            emailLabel.textColor = Style.Color.red
            return false
        } else if textField == passwordField {
            if let text = textField.text, text.characters.count > 0 {
                passwordLabel.textColor = Style.Color.black
                passwordField.resignFirstResponder()
                loginButtonTapped(loginButton)
                return true
            }
            passwordLabel.textColor = Style.Color.red
        }
        return false
    }
    
    func updateLabelColorsForTextField(_ textField: UITextField) {
        switch textField {
        case emailField: emailLabel.textColor = Style.Color.black
        case passwordField: passwordField.textColor = Style.Color.black
        default: break
        }
    }
}
