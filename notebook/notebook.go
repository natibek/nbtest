package notebook

//	{
//	    "metadata": {
//	        "kernel_info": {
//	            # if kernel_info is defined, its name field is required.
//	            "name": "the name of the kernel"
//	        },
//	        "language_info": {
//	            # if language_info is defined, its name field is required.
//	            "name": "the programming language of the kernel",
//	            "version": "the version of the language",
//	            "codemirror_mode": "The name of the codemirror mode to use [optional]",
//	        },
//	    },
//	    "nbformat": 4,
//	    "nbformat_minor": 0,
//	    "cells": [
//	        # list of cell dictionaries, see below
//	    ],
//	}
type notebookMetadata struct {
	KernelInfo   map[string]any    `json:"kernel_info"`
	LanguageInfo map[string]string `json:"language_info"`
}

type Notebook struct {
	Meatadata     notebookMetadata `json:"metadata"`
	NbFormat      uint             `json:"nbformat"`
	NbFormatMinor uint             `json:"nbformat_minor"`
	Cells         Cells            `json:"cells"`
}
