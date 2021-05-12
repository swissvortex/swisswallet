#!/bin/bash  
    
echo "Welcome to the swisswallet. You can currently generate, encrypt or decrypt a wallet." 
echo "Select one mode:"
echo "[1] Generate"
echo "[2] Encrypt"
echo "[3] Decrypt"
read -p "Mode: " mode  

echo ""
echo "All right. Now we need a password:"
read -sp "Type your password: " password  
echo ""
echo ""
echo "Just to double check. Please, type the password again:"
read -sp "Repeat your password: " repeatpassword  
echo ""
echo ""

if [ "$password" != "$repeatpassword" ]; then
    echo "Password does not match."
    echo ""
    exit
fi

if [ $mode == "1" ]; then
    echo "Almost there. We need a salt to proceed:"
    read -sp "Type your salt: " salt  
    echo ""
    echo ""
    echo "Just to double check. Please, type the salt again:"
    read -sp "Repeat your salt: " repeatsalt  
    echo ""
    echo ""
    if [ "$salt" != "$repeatsalt" ]; then
        echo "Salt does not match."
        echo ""
        exit
    fi
    echo "Generating a mnemonic from your password and salt..."
    go run main.go generate -p $password -s $salt
elif [ $mode == "2" ]; then
    echo "Please, type the mnemonic that you want to encrypt with the password:"
    read -p "Type your mnemonic (24 words): " mnemonic  
    echo ""
    echo "Encrypting the mnemonic with your password..."
    go run main.go encrypt -p $password -m "$mnemonic"
elif [ $mode == "3" ]; then
    echo "Please, type the address associated with the encrypted wallet (also printed when it was encrypted):"
    read -p "Type your address: " address  
    echo ""
    echo "Please, type the mnemonic that you want to decrypt with the password:"
    read -p "Type your encrypted mnemonic (24 words): " mnemonic  
    echo ""
    echo "Decrypting the mnemonic with your password and address..."
    go run main.go decrypt -p $password -m "$mnemonic" -a $address
else
  echo "Mode not supported"
fi

exit
