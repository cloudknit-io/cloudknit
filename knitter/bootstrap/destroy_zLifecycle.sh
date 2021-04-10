#!/bin/bash
# Copyright (C) 2020 CompuZest, Inc. - All Rights Reserved
#
# Unauthorized copying of this file, via any medium, is strictly prohibited
# Proprietary and confidential
#
# NOTICE: All information contained herein is, and remains the property of
# CompuZest, Inc. The intellectual and technical concepts contained herein are
# proprietary to CompuZest, Inc. and are protected by trade secret or copyright
# law. Dissemination of this information or reproduction of this material is
# strictly forbidden unless prior written permission is obtained from CompuZest, Inc.

main() {
    announcePhase() {
        echo ""
        echo ""
        echo "-------------------------------------"
        echo $1
        echo ""
        echo "-------------------------------------"   
    }

    checkForFailures() {
        if [ $? -ne 0 ]
        then
            echo ""
            echo "-------------------------------------"   
            read -p "Destroy phase has failed, type C to exit, any other key to continue" -n 1 -r
            echo ""

            if [[ $REPLY =~ ^[Cc]$ ]]
            then
                exit 1
            fi
        fi
    }

    echo "Please select the environment you wish to destroy:"
    select LOCATION in "dev-a" "dev-b" "sandbox" "demo"; do
        if [[ $LOCATION == "sandbox" || $LOCATION == "demo" ]]
        then
            readonly PARENT_DIRECTORY=aws
            break;
        else
            readonly PARENT_DIRECTORY=local
            break;
        fi
    done

    # Destroy
    announcePhase "Destroy"
    ./$PARENT_DIRECTORY/destroy.sh $LOCATION;
    checkForFailures

}

main
