const Layout = (props) => {
    return (
        <div className="px-14 md:px-40 lg:px-80 py-10 md:py-12 lg:py-15">
            {props.children}
        </div>
    )
}

export default Layout