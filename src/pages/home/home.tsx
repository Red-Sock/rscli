import cls from './home.module.css';

import {Header} from "../../sections/header/header";
import {ContentWrapper} from "../../sections/content/content";

import {Sidebar} from "../../sections/sidebar/sidebar";
import {Route, Routes} from "react-router-dom";
import {useState} from "react";
import {getResourceURLs} from "../../services/file-fetcher";

export function Home() {

    const [pageContent, setPageContent] = useState(`
# Uber Go Style Guide

- [Introduction](#introduction)
- [Guidelines](#guidelines)
  - [Pointers to Interfaces](#pointers-to-interfaces)
  - [Verify Interface Compliance](#verify-interface-compliance)
  - [Receivers and Interfaces](#receivers-and-interfaces)
  - [Zero-value Mutexes are Valid](#zero-value-mutexes-are-valid)
  - [Copy Slices and Maps at Boundaries](#copy-slices-and-maps-at-boundaries)
  - [Defer to Clean Up](#defer-to-clean-up)
  - [Channel Size is One or None](#channel-size-is-one-or-none)
  - [Start Enums at One](#start-enums-at-one)
  - [Use \`"time"\` to handle time](#use-time-to-handle-time)
  - [Errors](#errors)
    - [Error Types](#error-types)
    - [Error Wrapping](#error-wrapping)
    - [Error Naming](#error-naming)
    - [Handle Errors Once](#handle-errors-once)
  - [Handle Type Assertion Failures](#handle-type-assertion-failures)
  - [Don't Panic](#dont-panic)
  - [Use go.uber.org/atomic](#use-gouberorgatomic)
  - [Avoid Mutable Globals](#avoid-mutable-globals)
  - [Avoid Embedding Types in Public Structs](#avoid-embedding-types-in-public-structs)
  - [Avoid Using Built-In Names](#avoid-using-built-in-names)
  - [Avoid \`init()\`](#avoid-init)
  - [Exit in Main](#exit-in-main)
    - [Exit Once](#exit-once)
  - [Use field tags in marshaled structs](#use-field-tags-in-marshaled-structs)
  - [Don't fire-and-forget goroutines](#dont-fire-and-forget-goroutines)
    - [Wait for goroutines to exit](#wait-for-goroutines-to-exit)
    - [No goroutines in \`init()\`](#no-goroutines-in-init)
- [Performance](#performance)
  - [Prefer strconv over fmt](#prefer-strconv-over-fmt)
  - [Avoid string-to-byte conversion](#avoid-string-to-byte-conversion)
  - [Prefer Specifying Container Capacity](#prefer-specifying-container-capacity)
- [Style](#style)
  - [Avoid overly long lines](#avoid-overly-long-lines)
  - [Be Consistent](#be-consistent)
  - [Group Similar Declarations](#group-similar-declarations)
  - [Import Group Ordering](#import-group-ordering)
  - [Package Names](#package-names)
  - [Function Names](#function-names)
  - [Import Aliasing](#import-aliasing)
  - [Function Grouping and Ordering](#function-grouping-and-ordering)
  - [Reduce Nesting](#reduce-nesting)
  - [Unnecessary Else](#unnecessary-else)
  - [Top-level Variable Declarations](#top-level-variable-declarations)
  - [Prefix Unexported Globals with _](#prefix-unexported-globals-with-_)
  - [Embedding in Structs](#embedding-in-structs)
  - [Local Variable Declarations](#local-variable-declarations)
  - [nil is a valid slice](#nil-is-a-valid-slice)
  - [Reduce Scope of Variables](#reduce-scope-of-variables)
  - [Avoid Naked Parameters](#avoid-naked-parameters)
  - [Use Raw String Literals to Avoid Escaping](#use-raw-string-literals-to-avoid-escaping)
  - [Initializing Structs](#initializing-structs)
    - [Use Field Names to Initialize Structs](#use-field-names-to-initialize-structs)
    - [Omit Zero Value Fields in Structs](#omit-zero-value-fields-in-structs)
    - [Use \`var\` for Zero Value Structs](#use-var-for-zero-value-structs)
    - [Initializing Struct References](#initializing-struct-references)
  - [Initializing Maps](#initializing-maps)
  - [Format Strings outside Printf](#format-strings-outside-printf)
  - [Naming Printf-style Functions](#naming-printf-style-functions)
- [Patterns](#patterns)
  - [Test Tables](#test-tables)
  - [Functional Options](#functional-options)
- [Linting](#linting)


## Introduction

Styles are the conventions that govern our code. The term style is a bit of a
misnomer, since these conventions cover far more than just source file
formatting—gofmt handles that for us.

The goal of this guide is to manage this complexity by describing in detail the
Dos and Don'ts of writing Go code at Uber. These rules exist to keep the code
base manageable while still allowing engineers to use Go language features
productively.

This guide was originally created by [Prashant Varanasi](https://github.com/prashantv) and [Simon Newton](https://github.com/nomis52) as
a way to bring some colleagues up to speed with using Go. Over the years it has
been amended based on feedback from others.

This documents idiomatic conventions in Go code that we follow at Uber. A lot
of these are general guidelines for Go, while others extend upon external
resources:

1. [Effective Go](https://golang.org/doc/effective_go.html)
2. [Go Common Mistakes](https://github.com/golang/go/wiki/CommonMistakes)
3. [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

We aim for the code samples to be accurate for the two most recent minor versions
of Go [releases](https://go.dev/doc/devel/release).

All code should be error-free when run through \`golint\` and \`go vet\`. We
recommend setting up your editor to:

- Run \`goimports\` on save
- Run \`golint\` and \`go vet\` to check for errors


You can find information in editor support for Go tools here:
https://github.com/golang/go/wiki/IDEsAndTextEditorPlugins
`)

    function getContentViaLink(link: string) {
        getResourceURLs(link, setPageContent, (url: string)=> { window.location.replace(url)})
    }

    return (
        <>
            <div className={cls.headerWrap}>
                <Header/>
            </div>

            <div className={cls.Home}>

                <div className={cls.contentWrap}>
                    <Routes>
                        <Route path={"/*"} element={<ContentWrapper content={pageContent}/>}/>
                    </Routes>
                </div>

                <div className={cls.sideMenuWrap}>
                    <Sidebar setContentViaLink={getContentViaLink}/>
                </div>

            </div>
        </>
    )
}
